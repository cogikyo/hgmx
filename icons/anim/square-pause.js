import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll("line");
  animate(bars, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
