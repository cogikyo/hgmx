import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const horn = el.querySelector("g");
  animate(horn, { scale: [1, 0.9, 1], rotate: [0, -4, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
