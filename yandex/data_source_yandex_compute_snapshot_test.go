package yandex

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceComputeSnapshot(t *testing.T) {
	t.Parallel()

	diskName := acctest.RandomWithPrefix("tf-disk")
	snapshotName := acctest.RandomWithPrefix("tf-snap")
	label := acctest.RandomWithPrefix("label-value")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckComputeDiskDestroy,
			testAccCheckComputeSnapshotDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSnapshotConfig(diskName, snapshotName, label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.yandex_compute_snapshot.source",
						"name", snapshotName),
					resource.TestCheckResourceAttrSet("data.yandex_compute_snapshot.source",
						"id"),
					resource.TestCheckResourceAttrSet("data.yandex_compute_snapshot.source",
						"source_disk_id"),
					resource.TestCheckResourceAttr("data.yandex_compute_snapshot.source",
						"labels.test_label", label),
				),
			},
		},
	})
}

func testAccDataSourceSnapshotConfig(diskName, snapshotName, labelValue string) string {
	return fmt.Sprintf(`
data "yandex_compute_snapshot" "source" {
  snapshot_id = "${yandex_compute_snapshot.foobar.id}"
}

data "yandex_compute_image" "ubuntu" {
  family = "ubuntu-1804-lts"
}

resource "yandex_compute_disk" "foobar" {
  name     = "%s"
  image_id = "${data.yandex_compute_image.ubuntu.id}"
  size     = 4
}

resource "yandex_compute_snapshot" "foobar" {
  name           = "%s"
  source_disk_id = "${yandex_compute_disk.foobar.id}"

  labels = {
    test_label = "%s"
  }
}
`, diskName, snapshotName, labelValue)
}
