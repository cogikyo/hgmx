import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll("path");
  animate(bars, { scaleX: [1, 0.5, 1.1, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
