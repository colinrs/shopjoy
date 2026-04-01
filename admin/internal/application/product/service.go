package product

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"gorm.io/gorm"
)

type Service interface {
	CreateProduct(ctx context.Context, tenantID shared.TenantID, req CreateProductRequest) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, tenantID shared.TenantID, req UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, tenantID shared.TenantID, id int64) error
	GetProduct(ctx context.Context, tenantID shared.TenantID, id int64) (*ProductResponse, error)
	GetProductList(ctx context.Context, tenantID shared.TenantID, req QueryProductRequest) (*ProductListResponse, error)
	PutOnSale(ctx context.Context, tenantID shared.TenantID, id int64) (*ProductResponse, error)
	TakeOffSale(ctx context.Context, tenantID shared.TenantID, id int64) (*ProductResponse, error)
	UpdateStock(ctx context.Context, tenantID shared.TenantID, req UpdateStockRequest) error
	DeductStock(ctx context.Context, tenantID shared.TenantID, req DeductStockRequest) error
	BatchUpdateProduct(ctx context.Context, tenantID shared.TenantID, req BatchUpdateProductRequest) ([]int64, []BatchProductFail, error)
}

type service struct {
	db          *gorm.DB
	productRepo product.Repository
	idGen       snowflake.Snowflake
}

func NewService(db *gorm.DB, repo product.Repository, idGen snowflake.Snowflake) Service {
	return &service{
		db:          db,
		productRepo: repo,
		idGen:       idGen,
	}
}

func (s *service) CreateProduct(ctx context.Context, tenantID shared.TenantID, req CreateProductRequest) (*ProductResponse, error) {
	id, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, err
	}

	price := ToDomainMoney(req.Price, req.Currency)
	p, err := product.NewProduct(tenantID, req.Name, req.Description, price, req.CategoryID)
	if err != nil {
		return nil, err
	}
	p.ID = id

	p.CostPrice = ToDomainMoney(req.CostPrice, req.Currency)

	if err := s.productRepo.Create(ctx, s.db, p); err != nil {
		return nil, err
	}

	return FromDomainProduct(p), nil
}

func (s *service) UpdateProduct(ctx context.Context, tenantID shared.TenantID, req UpdateProductRequest) (*ProductResponse, error) {
	p, err := s.productRepo.FindByID(ctx, s.db, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	price := ToDomainMoney(req.Price, req.Currency)
	if err := p.UpdatePrice(price); err != nil {
		return nil, err
	}

	p.Name = req.Name
	p.Description = req.Description
	p.CategoryID = req.CategoryID
	p.SKU = req.SKU
	p.Brand = req.Brand
	p.Tags = req.Tags
	p.Images = req.Images
	p.IsMatrixProduct = req.IsMatrixProduct
	p.CostPrice = ToDomainMoney(req.CostPrice, req.Currency)
	p.HSCode = req.HSCode
	p.COO = req.COO
	p.Weight = req.Weight
	p.WeightUnit = req.WeightUnit
	p.Dimensions.Length = req.Length
	p.Dimensions.Width = req.Width
	p.Dimensions.Height = req.Height
	p.DangerousGoods = req.DangerousGoods

	if err := s.productRepo.Update(ctx, s.db, p); err != nil {
		return nil, err
	}

	return FromDomainProduct(p), nil
}

func (s *service) DeleteProduct(ctx context.Context, tenantID shared.TenantID, id int64) error {
	return s.productRepo.Delete(ctx, s.db, tenantID, id)
}

func (s *service) GetProduct(ctx context.Context, tenantID shared.TenantID, id int64) (*ProductResponse, error) {
	p, err := s.productRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return FromDomainProduct(p), nil
}

func (s *service) GetProductList(ctx context.Context, tenantID shared.TenantID, req QueryProductRequest) (*ProductListResponse, error) {
	query := product.Query{
		TenantID:   tenantID,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Status:     ParseStatus(req.Status),
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		MarketID:   req.MarketID,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	products, total, err := s.productRepo.FindList(ctx, s.db, query)
	if err != nil {
		return nil, err
	}

	resp := &ProductListResponse{
		List:     make([]*ProductResponse, len(products)),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	for i, p := range products {
		resp.List[i] = FromDomainProduct(p)
	}

	return resp, nil
}

func (s *service) PutOnSale(ctx context.Context, tenantID shared.TenantID, id int64) (*ProductResponse, error) {
	p, err := s.productRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := p.PutOnSale(); err != nil {
		return nil, err
	}

	if err := s.productRepo.Update(ctx, s.db, p); err != nil {
		return nil, err
	}

	return FromDomainProduct(p), nil
}

func (s *service) TakeOffSale(ctx context.Context, tenantID shared.TenantID, id int64) (*ProductResponse, error) {
	p, err := s.productRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := p.TakeOffSale(); err != nil {
		return nil, err
	}

	if err := s.productRepo.Update(ctx, s.db, p); err != nil {
		return nil, err
	}

	return FromDomainProduct(p), nil
}

func (s *service) UpdateStock(ctx context.Context, tenantID shared.TenantID, req UpdateStockRequest) error {
	p, err := s.productRepo.FindByID(ctx, s.db, tenantID, req.ID)
	if err != nil {
		return err
	}

	if err := p.UpdateStock(req.Quantity); err != nil {
		return err
	}

	return s.productRepo.Update(ctx, s.db, p)
}

func (s *service) DeductStock(ctx context.Context, tenantID shared.TenantID, req DeductStockRequest) error {
	return s.productRepo.UpdateStock(ctx, s.db, tenantID, req.ID, -req.Quantity)
}

// CreateProductWithTx 创建商品（带事务示例）
func (s *service) CreateProductWithTx(ctx context.Context, tenantID shared.TenantID, req CreateProductRequest) (*ProductResponse, error) {
	var result *product.Product

	err := s.db.Transaction(func(tx *gorm.DB) error {
		id, err := s.idGen.NextID(ctx)
		if err != nil {
			return err
		}

		price := ToDomainMoney(req.Price, req.Currency)
		p, err := product.NewProduct(tenantID, req.Name, req.Description, price, req.CategoryID)
		if err != nil {
			return err
		}
		p.ID = id

		p.CostPrice = ToDomainMoney(req.CostPrice, req.Currency)

		if err := s.productRepo.Create(ctx, tx, p); err != nil {
			return err
		}

		result = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return FromDomainProduct(result), nil
}

func (s *service) BatchUpdateProduct(ctx context.Context, tenantID shared.TenantID, req BatchUpdateProductRequest) ([]int64, []BatchProductFail, error) {
	var successIDs []int64
	var failed []BatchProductFail

	for _, productID := range req.ProductIDs {
		p, err := s.productRepo.FindByID(ctx, s.db, tenantID, productID)
		if err != nil {
			failed = append(failed, BatchProductFail{
				ProductID: productID,
				Code:      30012, // ErrProductNotFound code
				Message:   "商品不存在",
			})
			continue
		}

		// Apply updates
		if req.Fields.Price != nil {
			if err := p.UpdatePrice(product.Money{Amount: *req.Fields.Price, Currency: p.Price.Currency}); err != nil {
				failed = append(failed, BatchProductFail{
					ProductID: productID,
					Code:      30002, // ErrProductInvalidPrice code
					Message:   "商品价格必须大于0",
				})
				continue
			}
		}

		if req.Fields.Stock != nil {
			if err := p.UpdateStock(*req.Fields.Stock); err != nil {
				failed = append(failed, BatchProductFail{
					ProductID: productID,
					Code:      30008, // ErrProductNegativeStock code
					Message:   "库存不能为负数",
				})
				continue
			}
		}

		if req.Fields.Status != nil {
			// For status changes, we use the status transition methods
			switch *req.Fields.Status {
			case product.StatusOnSale:
				if err := p.PutOnSale(); err != nil {
					failed = append(failed, BatchProductFail{
						ProductID: productID,
						Code:      30006, // ErrProductInvalidStatusTransition code
						Message:   "无效的状态转换",
					})
					continue
				}
			case product.StatusOffSale:
				if err := p.TakeOffSale(); err != nil {
					failed = append(failed, BatchProductFail{
						ProductID: productID,
						Code:      30006,
						Message:   "无效的状态转换",
					})
					continue
				}
			default:
				failed = append(failed, BatchProductFail{
					ProductID: productID,
					Code:      30006,
					Message:   "无效的状态转换",
				})
				continue
			}
		}

		if req.Fields.CategoryID != nil {
			p.CategoryID = *req.Fields.CategoryID
		}

		if err := s.productRepo.Update(ctx, s.db, p); err != nil {
			failed = append(failed, BatchProductFail{
				ProductID: productID,
				Code:      20002, // ErrDatabase code
				Message:   "数据库错误",
			})
			continue
		}

		successIDs = append(successIDs, productID)
	}

	return successIDs, failed, nil
}
