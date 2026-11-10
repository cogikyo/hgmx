import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll('path[d="M10 13V7"], path[d="M14 13V7"]');
  animate(bars, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
