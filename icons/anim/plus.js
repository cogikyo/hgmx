import { animate, τmax } from "../motion.js";

export default function (el) {
  const bars = el.querySelectorAll("path");
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmax });
}
