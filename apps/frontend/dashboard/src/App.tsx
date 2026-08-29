import { Button } from '@klontong/ui';

/**
 * Kerangka dashboard (seller + admin). Fase 01 mengganti ini dengan
 * QueryClientProvider + RouterProvider + guard peran + MSW.
 */
export function App() {
  return (
    <main className="mx-auto flex min-h-screen max-w-2xl flex-col items-start gap-4 p-8">
      <h1 className="text-2xl font-semibold text-primary-700">Klontong Net — Dashboard</h1>
      <p className="text-muted-foreground">Kerangka Fase 00. Dasbor penjual &amp; panel admin.</p>
      <Button variant="secondary">Tombol contoh</Button>
    </main>
  );
}
