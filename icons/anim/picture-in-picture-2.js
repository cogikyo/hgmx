import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("path");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const inset = el.querySelector("rect");
  animate(inset, { scale: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
