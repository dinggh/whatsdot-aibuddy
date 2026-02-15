import * as React from "react";

import { cn } from "@/lib/utils";

const Card = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn(
        "rounded-2xl border border-[#E8E4DE] bg-white text-[#2D2A26] shadow-[0_2px_8px_rgba(26,25,24,0.03)]",
        className
      )}
      {...props}
    />
  )
);
Card.displayName = "Card";

export { Card };
