import { animate, τmax } from "../motion.js";

export default function (el) {
  const bar = el.querySelector("path");
  animate(bar, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmax });
}
