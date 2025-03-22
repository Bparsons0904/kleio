# Kleio Style Guidelines

This document outlines the style guidelines for the Kleio project to ensure consistency across all UI components and pages.

## Design Principles

- **Clarity**: UI elements should be easily understandable at a glance
- **Consistency**: Similar elements should look and behave similarly across the application
- **Responsiveness**: All components should work well across desktop, tablet, and mobile devices
- **Accessibility**: The application should be usable by people with diverse abilities

## Color Palette

### Primary Colors

- Primary Blue: `#3273dc` (Used for buttons, links, and primary actions)
- Darker Blue: Darken primary blue by 5% for hover states (`darken(#3273dc, 5%)`)

### Text Colors

- Headings: `#333` (Dark gray)
- Subheadings: `#444` (Medium-dark gray)
- Body Text: `#555` (Medium gray)
- Secondary Text: `#666` (Medium-light gray)
- Disabled Text: `#777` (Light gray)

### Background Colors

- Page Background: White or very light gray
- Card Background: White
- Card Header: `#f8f9fa` (Very light gray)
- Card Footer: White

### Feedback Colors

- Success: `#2f855a` (Green)
- Success Background: `#f0fff4` (Light green)
- Error: `#c53030` (Red)
- Error Background: `#fff5f5` (Light red)
- Info: `#2b6cb0` (Blue)
- Info Background: `#ebf8ff` (Light blue)

## Typography

### Font Hierarchy

- Primary Heading (h1): `2rem`, `700` weight
- Secondary Heading (h2): `1.75rem`, `600` weight
- Card Heading: `1.3rem`, `600` weight
- Body Text: `1rem`, `400` weight
- Small Text: `0.875rem`

### Line Heights

- Headings: `1.2`
- Body Text: `1.5`

## Spacing

### Margins and Padding

- Container Padding: `2rem`
- Card Padding: `1.5rem`
- Section Margin: `2rem` bottom
- Component Margin: `1rem` bottom
- Form Group Margin: `1rem` bottom
- Button Padding: `0.75rem 1.5rem` (vertical, horizontal)

### Grid System

- Use CSS Grid for page layouts
- Card Grid: `repeat(auto-fill, minmax(300px, 1fr))`
- Grid Gap: `1.5rem`

## Components

### Containers

- Max width: `1200px` for main content
- Centered with `margin: 0 auto`
- Responsive padding: `2rem` on desktop, `1.5rem` on mobile

### Cards

```scss
.card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition:
    transform 0.2s,
    box-shadow 0.2s;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
  }
}

.cardHeader {
  padding: 1.25rem 1.5rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #eee;
}

.cardBody {
  padding: 1.5rem;
  flex-grow: 1;
}

.cardFooter {
  padding: 1rem 1.5rem;
  border-top: 1px solid #eee;
}
```

### Buttons

```scss
.button {
  padding: 0.75rem 1.5rem;
  background-color: #3273dc;
  color: white;
  font-weight: 500;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition:
    background-color 0.2s,
    opacity 0.2s;

  &:hover:not(:disabled) {
    background-color: darken(#3273dc, 5%);
  }

  &:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(50, 115, 220, 0.25);
  }

  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
}
```

### Form Elements

```scss
.input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;

  &:focus {
    outline: none;
    border-color: #3273dc;
    box-shadow: 0 0 0 2px rgba(50, 115, 220, 0.25);
  }
}

.label {
  display: block;
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: #444;
}
```

### Messages

```scss
.message {
  margin-bottom: 1rem;
  padding: 0.75rem;
  border-radius: 4px;

  &.success {
    background-color: #f0fff4;
    color: #2f855a;
    border-left: 4px solid #2f855a;
  }

  &.error {
    background-color: #fff5f5;
    color: #c53030;
    border-left: 4px solid #c53030;
  }

  &.info {
    background-color: #ebf8ff;
    color: #2b6cb0;
    border-left: 4px solid #3182ce;
  }
}
```

## Responsive Breakpoints

- Mobile: `480px`
- Tablet: `768px`
- Desktop: `1024px`
- Large Desktop: `1280px`

```scss
// Mobile styles
@media (max-width: 480px) {
  // Mobile-specific styles
}

// Tablet styles
@media (max-width: 768px) {
  // Tablet-specific styles
}

// Desktop styles
@media (max-width: 1024px) {
  // Desktop-specific styles
}
```

## CSS Architecture

1. Use CSS Modules with SolidJS for component-scoped styles
2. Maintain consistency by importing shared variables from a central file
3. Follow a naming convention for class names:
   - Use camelCase for class names
   - Be descriptive but concise
   - Use semantic names that describe the purpose, not the appearance

## Animation Guidelines

- Keep animations subtle and purposeful
- Standard transition duration: `0.2s`
- Use transitions for hover states and UI feedback
- Avoid animations that could cause motion sickness or distraction

## Best Practices

1. **Mobile-First Approach**: Design for mobile first, then enhance for larger screens
2. **Consistent Spacing**: Use the defined spacing values consistently
3. **Semantic HTML**: Use appropriate HTML elements for their intended purpose
4. **Progressive Enhancement**: Ensure basic functionality works without JavaScript or advanced CSS
5. **Code Comments**: Add comments to explain complex or non-obvious styling decisions

## File Structure

```
src/
├── components/
│   ├── ComponentName/
│   │   ├── ComponentName.tsx
│   │   └── ComponentName.module.scss
├── pages/
│   ├── PageName/
│   │   ├── PageName.tsx
│   │   └── PageName.module.scss
├── styles/
│   ├── _variables.scss  (for shared variables)
│   └── global.scss      (for global styles)
```

By following these guidelines, we'll maintain a consistent and professional look and feel throughout the Kleio application, making it easier to develop new features while ensuring a cohesive user experience.
