import { animate, τmid } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll("rect");
  animate(bars, { scaleY: [0, 1.15, 1] }, { duration: τmid });
}
