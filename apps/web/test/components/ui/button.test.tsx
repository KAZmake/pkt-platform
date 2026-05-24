import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Button } from '@/components/ui/button';

describe('Button', () => {
  it('renders children', () => {
    render(<Button>Нажми меня</Button>);
    expect(screen.getByRole('button', { name: 'Нажми меня' })).toBeInTheDocument();
  });

  it('is enabled by default', () => {
    render(<Button>OK</Button>);
    expect(screen.getByRole('button')).not.toBeDisabled();
  });

  it('is disabled when disabled prop is set', () => {
    render(<Button disabled>OK</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('is disabled when isLoading is true', () => {
    render(<Button isLoading>OK</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('shows spinner when isLoading', () => {
    const { container } = render(<Button isLoading>OK</Button>);
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toBeInTheDocument();
  });

  it('no spinner when not loading', () => {
    const { container } = render(<Button>OK</Button>);
    expect(container.querySelector('.animate-spin')).not.toBeInTheDocument();
  });

  it('applies primary variant class by default', () => {
    render(<Button>OK</Button>);
    const btn = screen.getByRole('button');
    expect(btn.className).toContain('bg-brand-green');
  });

  it('applies destructive variant classes', () => {
    render(<Button variant="destructive">Delete</Button>);
    const btn = screen.getByRole('button');
    expect(btn.className).toContain('bg-red-600');
  });

  it('applies outline variant classes', () => {
    render(<Button variant="outline">Outline</Button>);
    const btn = screen.getByRole('button');
    expect(btn.className).toContain('border-brand-green');
  });

  it('applies sm size classes', () => {
    render(<Button size="sm">Small</Button>);
    const btn = screen.getByRole('button');
    expect(btn.className).toContain('h-8');
  });

  it('applies lg size classes', () => {
    render(<Button size="lg">Large</Button>);
    const btn = screen.getByRole('button');
    expect(btn.className).toContain('h-12');
  });

  it('merges custom className', () => {
    render(<Button className="custom-class">OK</Button>);
    expect(screen.getByRole('button').className).toContain('custom-class');
  });

  it('passes onClick handler', async () => {
    let clicked = false;
    render(
      <Button
        onClick={() => {
          clicked = true;
        }}
      >
        Click
      </Button>,
    );
    screen.getByRole('button').click();
    expect(clicked).toBe(true);
  });
});
