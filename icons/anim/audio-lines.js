import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll("path");
  animate(bars, { scaleY: [1, 0.6, 1.2, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
