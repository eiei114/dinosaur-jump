package _interface

import (
	"encoding/json"
	"example.com/application/service"
	"example.com/interface/request"
	"example.com/interface/response"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: *userService}
}

func (u *UserHandler) UserCreateHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var requestData request.UserCreateRequest
		if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return err
		}

		ctx := req.Context()
		authToken, err := u.userService.Add(ctx, requestData.Name)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return err
		}

		responseData := &response.UserCreateResponse{Token: authToken}
		responseBytes, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
		return nil
	}
}

// UserGetHandle retrieves user information based on auth_token
func (u *UserHandler) UserGetHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var requestData request.UserGetRequest

		// Decode the request body to get auth_token
		if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return err
		}

		ctx := req.Context()

		// Retrieve user by auth token
		user, err := u.userService.GetUserByAuthToken(ctx, requestData.Token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		// Prepare the response using UserGetResponse struct
		responseData := &response.UserGetResponse{
			Id:        user.Id,
			Name:      user.Name,
			HighScore: user.HighScore,
		}

		respBytes, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)

		return nil
	}
}

// MoveHandle プレイヤー移動同期
func (u *UserHandler) MoveHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		w.Write([]byte("MoveHandle triggered"))
		return nil
	}
}

func (u *UserHandler) UserRankingGetHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// UserServiceからランキングを取得
		userRankings, err := u.userService.GetUserRanking(req.Context())
		if err != nil {
			http.Error(w, "Failed to get user rankings", http.StatusInternalServerError)
			return err
		}

		// UserRankingからUserRankingResponseに変換
		var responseSlice []response.UserRankingResponse
		for _, ranking := range userRankings {
			r := response.UserRankingResponse{
				Name:      ranking.Name,
				HighScore: ranking.HighScore,
			}
			responseSlice = append(responseSlice, r)
		}

		// ヘッダーを設定してJSONとしてレスポンスを返すことを示す
		w.Header().Set("Content-Type", "application/json")

		// ランキングをJSONとしてエンコードしてレスポンスに書き込む
		if err := json.NewEncoder(w).Encode(responseSlice); err != nil {
			http.Error(w, "Failed to encode user rankings", http.StatusInternalServerError)
			return err
		}

		return nil
	}
}

// DestroyHandle プレイヤーゲームオーバー
func (u *UserHandler) DestroyHandle() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		w.Write([]byte("DestroyHandle triggered"))
		return nil
	}
}
