import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cross = el.querySelector("path");
  animate(cross, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}
