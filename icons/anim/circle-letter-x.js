import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const stroke = el.querySelector('path[d="m9 16 6-8"]');
  animate(stroke, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const cross = el.querySelector('path[d="m9 8 6 8"]');
  animate(cross, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
