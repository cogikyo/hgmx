import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const s = el.querySelector('path[d^="M16 8"]');
  animate(s, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const bar = el.querySelector('path[d="M12 18V6"]');
  animate(bar, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
