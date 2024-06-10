package handlers

import (
	"log"
	"net/http"

	common "github.com/millukii/commons"
	pb "github.com/millukii/commons/api"
	errors "github.com/millukii/commons/errors"
	"github.com/millukii/openmarket-gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHttpHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway: gateway}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux){

		mux.HandleFunc("POST /api/customer/{customerID}/orders", h.HandleCreateOrder )

}



func (h *handler) HandleCreateOrder(w http.ResponseWriter, r  *http.Request){

	customerID := r.URL.Query().Get("customerID")
	
	var items []*pb.ItemsWithQuantity

	if err := common.ReadJSON(r, &items); err != nil {
		log.Println("error reading json", err)
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return 
	}	

	if err := validateItems(items); err != nil{
		log.Println("error at items validation", err)
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	order, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerOd: customerID,
		Items: items,
	})

	rStatus := status.Convert(err)
	
	if rStatus != nil{
		if rStatus.Code() != codes.InvalidArgument{
					common.WriteError(w, http.StatusBadRequest, rStatus.Message())
		}
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.WriteJSON(w, http.StatusOK, order)
}

func validateItems(items []*pb.ItemsWithQuantity) error{
	if (len(items) ==0){
		return errors.ErrNoItems
	}

	for _,i :=range items{
		if i.ID == ""{
				return errors.ErrNoId
		}
		if i.Quantity <= 0{
				return errors.ErrInvalidQuantity
		}
	}
	mergedItems := mergeItemsQuantities(items)

	log.Println("merged items ", mergedItems)
	//validate with stock service
	return nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity{
	merged := make([]*pb.ItemsWithQuantity,0)

	for _, item := range items{
		found := false

		for _,finalItem := range merged {
			if finalItem.ID == item.ID{
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
			if !found{
				merged = append(merged, item)
			}
		}
	}

	return merged
}