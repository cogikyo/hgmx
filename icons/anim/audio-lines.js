import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll("path");
  animate(bars, { scaleY: [1, 0.4, 1.3, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
