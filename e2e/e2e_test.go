package e2e

import (
	"banner-service/internal/app"
	controller "banner-service/internal/controller/http"
	"banner-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"testing"
)

const (
	addr          = "http://localhost:8000"
	testFeatureID = 123456789
	testContent   = `{"hello": "world"}`

	newTestContent   = `{"hello": "new world"}`
	newTestFeatureID = 987654321
)

var (
	testTagIDs    = []uint64{9000100, 9000200, 9000300, 9000400}
	newTestTagIds = []uint64{12345, 6789}
)

func Setup() {
	go app.Start()

	const truncateQuery = `
		truncate banner, banner_feature_tag, banner_version, role_endpoints, users
	`

	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432")
	if err != nil {
		panic(err)
	}
	_, err = pool.Exec(context.Background(), truncateQuery)
	if err != nil {
		panic(err)
	}

	const addResourcesQuery = `
		insert into role_endpoints (role, resource)
		values
    		('admin', '*'),
    		('user', 'GET /user_banner')`

	_, err = pool.Exec(context.Background(), addResourcesQuery)
	if err != nil {
		panic(err)
	}
}

type testClient struct {
	resty *resty.Client
}

func (c testClient) SignUp(user models.User) (*resty.Response, error) {
	return c.resty.R().SetBody(user).Post(addr + "/sign-up")
}

func (c testClient) SignIn(user models.User) (*resty.Response, error) {
	return c.resty.R().SetBody(user).Post(addr + "/sign-in")
}

func (c testClient) CreateBanner(banner controller.CreateDTO, token string) (*resty.Response, error) {
	return c.resty.R().
		SetBody(banner).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Post(addr + "/banner")
}

func (c testClient) GetBanner(tagID, featureID uint64, token string) (*resty.Response, error) {
	return c.resty.R().SetQueryParams(map[string]string{
		"tag_id":     fmt.Sprint(tagID),
		"feature_id": fmt.Sprint(featureID),
	}).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Get(addr + "/user_banner")
}

func (c testClient) PatchBanner(bannerId uint64, banner controller.CreateDTO, token string) (*resty.Response, error) {
	return c.resty.R().
		SetBody(banner).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Patch(fmt.Sprintf("%s/banner/%d", addr, bannerId))
}

func (c testClient) DeleteBanner(bannerId uint64, token string) (*resty.Response, error) {
	return c.resty.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Delete(fmt.Sprintf("%s/banner/%d", addr, bannerId))
}

func (c testClient) GetFilteredBanners(filterBanner models.FilterBanner, token string) (*resty.Response, error) {
	return c.resty.R().SetQueryParams(map[string]string{
		"tag_id":     fmt.Sprint(filterBanner.TagId),
		"feature_id": fmt.Sprint(filterBanner.FeatureId),
		"limit":      fmt.Sprint(filterBanner.Limit),
		"offset":     fmt.Sprint(filterBanner.Offset),
	}).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Get(addr + "/banner")
}

func TestBasicScenario(t *testing.T) {
	Setup()

	client := testClient{resty.New()}

	var userToken string
	t.Run("auth", func(t *testing.T) {
		resp, err := client.SignUp(models.User{
			Username: "username",
			Password: "password",
			Role:     "user",
		})
		assert.NoError(t, err)
		userToken = string(resp.Body())

		assert.Equal(t, http.StatusCreated, resp.StatusCode())
		if !assert.NotEmpty(t, userToken) {
			return
		}
	})

	var adminToken string
	t.Run("auth", func(t *testing.T) {
		resp, err := client.SignUp(models.User{
			Username: "qwerty",
			Password: "12345678",
			Role:     "admin",
		})
		assert.NoError(t, err)
		adminToken = string(resp.Body())

		assert.Equal(t, http.StatusCreated, resp.StatusCode())
		if !assert.NotEmpty(t, adminToken) {
			return
		}
	})

	var bannerId uint64
	t.Run("create banner", func(t *testing.T) {
		bannerDTO := controller.CreateDTO{
			FeatureId: testFeatureID,
			TagIds:    testTagIDs,
			Content:   json.RawMessage(testContent),
			IsActive:  false,
		}

		resp, err := client.CreateBanner(bannerDTO, adminToken)
		if err != nil {
			t.Fatal(err)
		}
		bannerId, err = strconv.ParseUint(string(resp.Body()), 10, 64)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode())
	})

	t.Run("get banner", func(t *testing.T) {
		resp, err := client.GetBanner(testTagIDs[0], testFeatureID, userToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode())
	})

	t.Run("get banner", func(t *testing.T) {
		resp, err := client.GetBanner(testTagIDs[0], testFeatureID, adminToken)
		assert.NoError(t, err)
		content := string(resp.Body())
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.Equal(t, testContent, content)
	})

	t.Run("get banner", func(t *testing.T) {
		resp, err := client.GetBanner(testTagIDs[0], testFeatureID, "fake token")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})

	t.Run("patch banner", func(t *testing.T) {
		banner := controller.CreateDTO{
			FeatureId: newTestFeatureID,
			TagIds:    newTestTagIds,
			Content:   json.RawMessage(newTestContent),
			IsActive:  true}
		resp, err := client.PatchBanner(bannerId, banner, adminToken)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("get banner", func(t *testing.T) {
		resp, err := client.GetBanner(newTestTagIds[0], newTestFeatureID, userToken)
		assert.NoError(t, err)
		content := string(resp.Body())
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.Equal(t, newTestContent, content)
	})

	t.Run("get filtered banner", func(t *testing.T) {
		resp, err := client.GetFilteredBanners(models.FilterBanner{
			FeatureId: testFeatureID,
			TagId:     testTagIDs[0],
			Limit:     1,
			Offset:    0,
		}, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("delete banner", func(t *testing.T) {
		resp, err := client.DeleteBanner(bannerId, adminToken)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	})

}

func Test(t *testing.T) {
	Setup()

	client := testClient{resty.New()}

	resp, err := client.SignUp(models.User{
		Username: "qwerty",
		Password: "12345678",
		Role:     "admin",
	})
	if err != nil {
		t.Fatal(err)
	}

	token := string(resp.Body())

	eg := errgroup.Group{}

	for i := uint64(0); i < 128; i++ {
		eg.Go(func() error {
			_, err := client.CreateBanner(controller.CreateDTO{
				FeatureId: i,
				TagIds:    []uint64{i},
				Content:   json.RawMessage(fmt.Sprintf(`{"bannerNo": %d}`, i)),
				IsActive:  true,
			}, token)
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	eg = errgroup.Group{}

	for i := uint64(0); i < uint64(256); i++ {
		i := i
		eg.Go(func() error {
			n := i % 16
			_, err := client.GetBanner(n, n, token)
			if err != nil {
				fmt.Println(err)
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestStressWithCache(t *testing.T) {
	Setup()

	limit := 1000
	rateLimiter := rate.NewLimiter(rate.Limit(limit), limit)

	client := testClient{resty: resty.New().SetRetryCount(5)}

	resp, err := client.SignUp(models.User{
		Username: "qwerty",
		Password: "12345678",
		Role:     "admin",
	})
	if err != nil {
		t.Fatal(err)
	}

	token := string(resp.Body())

	eg := errgroup.Group{}

	const bannerCount = 200

	for i := uint64(0); i < bannerCount; i++ {
		eg.Go(func() error {
			_, err := client.CreateBanner(controller.CreateDTO{
				FeatureId: i,
				TagIds:    []uint64{i},
				Content:   json.RawMessage(fmt.Sprintf(`{"bannerNo": %d}`, i)),
				IsActive:  true,
			}, token)
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	eg = errgroup.Group{}

	for i := uint64(0); ; i++ {
		if err := rateLimiter.Wait(context.Background()); err != nil {
			fmt.Println("rate limiter:", err)
		}

		i := i
		rateLimiter.Reserve()
		go func(i uint64) {
			defer rateLimiter.Allow()
			client := testClient{resty: resty.New()}
			n := i % bannerCount
			_, _ = client.GetBanner(n, n, token)
		}(i)
	}
}
