import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { App } from './App';

describe('App dashboard', () => {
  it('merender judul dan tombol contoh', () => {
    render(<App />);
    expect(screen.getByRole('heading', { name: /dashboard/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /tombol contoh/i })).toBeInTheDocument();
  });
});
