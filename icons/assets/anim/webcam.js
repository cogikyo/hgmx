import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector('circle[r="8"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const lens = el.querySelector('circle[r="3"]');
  animate(lens, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
