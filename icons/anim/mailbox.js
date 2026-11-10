import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M22 17"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flag = el.querySelector("polyline");
  animate(flag, { rotate: [0, -30, 10, 0] }, { duration: τmid, delay: τmin });
}
