import { cn } from '@/lib/utils';

type SiteLogoIconProps = {
  className?: string;
  title?: string;
};

export function SiteLogoIcon({
  className,
  title = 'Unfollow Tracker logo',
}: SiteLogoIconProps) {
  return (
    <img src="/logo.png" alt={title} className={className} />
  );
}

type SiteLogoProps = {
  className?: string;
  iconClassName?: string;
  textClassName?: string;
  withText?: boolean;
};

export function SiteLogo({
  className,
  iconClassName,
  textClassName,
  withText = true,
}: SiteLogoProps) {
  return (
    <div className={cn('flex items-center gap-3', className)}>
      <SiteLogoIcon className={cn('h-10 w-10 shrink-0', iconClassName)} />
      {withText && (
        <span className={cn('text-xl font-heading font-semibold text-text-primary', textClassName)}>
          Unfollow Tracker
        </span>
      )}
    </div>
  );
}
