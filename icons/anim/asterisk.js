import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll("path");
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06) });
}
