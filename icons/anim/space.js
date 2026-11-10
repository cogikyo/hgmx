import { animate, τmid } from "../motion.js";

export default function (el) {
  const bar = el.querySelector("path");
  animate(bar, { translateY: [0, 2, 0] }, { duration: τmid });
}
