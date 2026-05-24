import { ActivityIndicator, Pressable, Text, type PressableProps } from 'react-native';

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger';
type Size = 'sm' | 'md' | 'lg';

interface ButtonProps extends PressableProps {
  variant?: Variant;
  size?: Size;
  loading?: boolean;
  label: string;
}

const variantStyles: Record<Variant, { bg: string; text: string }> = {
  primary: { bg: 'bg-brand-green-500 active:bg-brand-green-600', text: 'text-white' },
  secondary: { bg: 'bg-brand-gold-500 active:bg-brand-gold-dark', text: 'text-white' },
  ghost: { bg: 'bg-transparent border border-brand-green-500', text: 'text-brand-green-500' },
  danger: { bg: 'bg-red-500 active:bg-red-600', text: 'text-white' },
};

const sizeStyles: Record<Size, { container: string; text: string }> = {
  sm: { container: 'px-3 py-1.5 rounded-lg', text: 'text-sm font-medium' },
  md: { container: 'px-4 py-2.5 rounded-xl', text: 'text-base font-semibold' },
  lg: { container: 'px-6 py-3.5 rounded-xl', text: 'text-lg font-semibold' },
};

export function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  label,
  disabled,
  className,
  ...props
}: ButtonProps) {
  const isDisabled = disabled || loading;
  const v = variantStyles[variant];
  const s = sizeStyles[size];

  return (
    <Pressable
      className={`items-center justify-center flex-row gap-2 ${v.bg} ${s.container} ${isDisabled ? 'opacity-50' : ''} ${className ?? ''}`}
      disabled={isDisabled}
      {...props}
    >
      {loading && <ActivityIndicator size="small" color="white" />}
      <Text className={`${v.text} ${s.text}`}>{label}</Text>
    </Pressable>
  );
}
