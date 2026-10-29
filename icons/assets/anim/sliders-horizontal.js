import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const outer = el.querySelectorAll('path[d="M14 3v4"], path[d="M16 17v4"]');
  animate(outer, { translateX: [0, -4, 1, 0] }, { duration: 0.3, delay: stagger(0.12) });

  const middle = el.querySelector('path[d="M8 10v4"]');
  animate(middle, { translateX: [0, 4, -1, 0] }, { duration: τmax });
}
