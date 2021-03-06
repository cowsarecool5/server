package dashboard

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestDashboard(t *testing.T) {

	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1)
	db.User(2)

	resolver := NewResolverForDashboard(db.DB)

	// initial no dashboards
	user1 := fake.User(1)
	user2 := fake.User(2)
	dashboards, err := resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Empty(t, dashboards)

	dashboard, err := resolver.CreateDashboard(user1, "cool dashboard")
	require.NoError(t, err)
	expectAdded := &gqlmodel.Dashboard{
		ID:     1,
		Name:   "cool dashboard",
		Ranges: []*gqlmodel.NamedDateRange{},
		Items:  []*gqlmodel.DashboardEntry{},
	}
	require.Equal(t, expectAdded, dashboard)
	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Dashboard{
		{
			ID:     1,
			Name:   "cool dashboard",
			Ranges: []*gqlmodel.NamedDateRange{},
			Items:  []*gqlmodel.DashboardEntry{},
		},
	}, dashboards)
	dashboards, err = resolver.Dashboards(user2)
	require.NoError(t, err)
	require.Empty(t, dashboards)

	dashboard, err = resolver.UpdateDashboard(user2, 1, "kek")
	require.EqualError(t, err, "dashboard does not exist")
	require.Nil(t, dashboard)

	dashboard, err = resolver.UpdateDashboard(user2, 55, "kek")
	require.EqualError(t, err, "dashboard does not exist")
	require.Nil(t, dashboard)

	dashboard, err = resolver.UpdateDashboard(user1, 1, "kek")
	require.NoError(t, err)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Dashboard{
		{
			ID:     1,
			Name:   "kek",
			Ranges: []*gqlmodel.NamedDateRange{},
			Items:  []*gqlmodel.DashboardEntry{},
		},
	}, dashboards)

	dashboard, err = resolver.RemoveDashboard(user2, 1)
	require.EqualError(t, err, "dashboard does not exist")
	require.Nil(t, dashboard)

	dashboard, err = resolver.RemoveDashboard(user2, 55)
	require.EqualError(t, err, "dashboard does not exist")
	require.Nil(t, dashboard)

	dashboard, err = resolver.RemoveDashboard(user1, 1)
	require.NoError(t, err)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Empty(t, dashboards)
}

func TestDeleteDashboard(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	dbOne := db.User(1).Dashboard("test")
	dbOne.Range("r")
	dbOne.Entry("e")
	dbTwo := db.User(2).Dashboard("test")
	dbTwo.Range("r")
	dbTwo.Entry("e")

	resolver := NewResolverForDashboard(db.DB)
	_, err := resolver.RemoveDashboard(fake.User(dbOne.User.ID), dbOne.Dashboard.ID)
	require.NoError(t, err)

	dbOne.AssertExists(false).AssertHasEntry("e", false).AssertHasRange("r", false)
	dbTwo.AssertExists(true).AssertHasEntry("e", true).AssertHasRange("r", true)
}
