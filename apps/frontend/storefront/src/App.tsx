import { Button } from '@klontong/ui';

/**
 * Kerangka storefront (guest + buyer). Fase 01 mengganti ini dengan
 * QueryClientProvider + RouterProvider + MSW.
 */
export function App() {
  return (
    <main className="mx-auto flex min-h-screen max-w-2xl flex-col items-start gap-4 p-8">
      <h1 className="text-2xl font-semibold text-primary-700">Klontong Net — Storefront</h1>
      <p className="text-muted-foreground">
        Kerangka Fase 00. Belanja untuk pembeli (guest &amp; login).
      </p>
      <Button>Tombol contoh</Button>
    </main>
  );
}
