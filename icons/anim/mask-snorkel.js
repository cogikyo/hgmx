import { animate, stagger } from "../motion.js";

export default function (el) {
  const bubbles = el.querySelectorAll('circle, path[d="M3 17.5h.01"]');
  animate(bubbles, { translateY: [0, -3, 0.5, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
