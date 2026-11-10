import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector('circle[r="10"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const bowl = el.querySelector('circle[r="2"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: τmid });

  const tail = el.querySelector("path");
  animate(tail, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
