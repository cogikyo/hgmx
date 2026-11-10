import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const outer = el.querySelectorAll('path[d="M3 14h4"], path[d="M17 16h4"]');
  animate(outer, { translateY: [0, -4, 1, 0] }, { duration: 0.3, delay: stagger(0.12) });

  const middle = el.querySelector('path[d="M10 8h4"]');
  animate(middle, { translateY: [0, 4, -1, 0] }, { duration: τmax });
}
