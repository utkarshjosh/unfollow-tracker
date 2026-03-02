import type { ReactElement } from 'react';

interface HowToStep {
  name: string;
  text: string;
  image?: string;
  url?: string;
}

interface HowToSchemaProps {
  name: string;
  description: string;
  steps: HowToStep[];
  totalTime?: string; // ISO 8601 duration format (e.g., "PT5M")
  estimatedCost?: {
    currency: string;
    value: string;
  };
  image?: string;
}

interface HowToStepSchema {
  "@type": "HowToStep";
  name: string;
  text: string;
  image?: string;
  url?: string;
}

interface EstimatedCostSchema {
  "@type": "MonetaryAmount";
  currency: string;
  value: string;
}

interface HowToSchemaData {
  "@context": "https://schema.org";
  "@type": "HowTo";
  name: string;
  description: string;
  totalTime?: string;
  estimatedCost?: EstimatedCostSchema;
  image?: string;
  step: HowToStepSchema[];
}

/**
 * HowToSchema Component
 *
 * Renders JSON-LD HowTo schema markup for rich snippets in Google search results.
 * Follows schema.org/HowTo specification.
 *
 * @example
 * // How to track unfollowers
 * <HowToSchema
 *   name="How to Track Instagram Unfollowers"
 *   description="Step-by-step guide to tracking unfollowers..."
 *   totalTime="PT5M"
 *   estimatedCost={{ currency: "USD", value: "0" }}
 *   steps={[
 *     { name: "Create Account", text: "Sign up for a free Unfollow Tracker account..." },
 *     { name: "Connect Instagram", text: "Securely connect your Instagram account..." },
 *     { name: "View Insights", text: "Access your unfollower data..." }
 *   ]}
 * />
 */
function HowToSchema({
  name,
  description,
  steps,
  totalTime,
  estimatedCost,
  image,
}: HowToSchemaProps): ReactElement {
  const schemaData: HowToSchemaData = {
    "@context": "https://schema.org",
    "@type": "HowTo",
    name,
    description,
    step: steps.map((step): HowToStepSchema => ({
      "@type": "HowToStep",
      name: step.name,
      text: step.text,
      ...(step.image && { image: step.image }),
      ...(step.url && { url: step.url }),
    })),
  };

  if (totalTime) {
    schemaData.totalTime = totalTime;
  }

  if (estimatedCost) {
    schemaData.estimatedCost = {
      "@type": "MonetaryAmount",
      currency: estimatedCost.currency,
      value: estimatedCost.value,
    };
  }

  if (image) {
    schemaData.image = image;
  }

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{
        __html: JSON.stringify(schemaData),
      }}
    />
  );
}

export type { HowToStep, HowToSchemaProps };
export default HowToSchema;
