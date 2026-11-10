import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const body = el.querySelector('path[d^="M10 16"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const bars = el.querySelectorAll('path[d="M8 12h4"], path[d="M8 16h7"]');
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
