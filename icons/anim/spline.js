import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const start = el.querySelector('circle[cx="5"]');
  animate(start, { scale: [1, 1.25, 1] }, { duration: τmin });

  const curve = el.querySelector("path");
  animate(curve, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const end = el.querySelector('circle[cx="19"]');
  animate(end, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
