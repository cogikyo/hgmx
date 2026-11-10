import { animate, stagger } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll("rect");
  animate(bars, { scaleX: [0.7, 1.1, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
