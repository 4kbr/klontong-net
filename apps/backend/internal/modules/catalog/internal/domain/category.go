package domain

// TODO:
//   type Category struct { ID; ParentID *uuid.UUID; Name; Slug; IconURL;
//                          Position int; IsActive; CreatedAt; UpdatedAt }
//   func (c *Category) IsRoot() bool
//
// Berjenjang. Untuk menampilkan pohon tanpa recursive CTE di setiap request,
// simpan `path` (materialized path atau ltree) dan perbarui saat parent berubah.
// Kategori jarang berubah dan sering dibaca — kandidat kuat untuk di-cache.
