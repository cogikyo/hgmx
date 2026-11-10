import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll('path[d="M6 9v6"], path[d="M10 6v12"], path[d="M14 6v12"], path[d="M18 9v6"]');
  animate(bars, { scaleY: [1, 0.3, 1.3, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
