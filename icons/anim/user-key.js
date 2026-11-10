import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M3 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const key = el.querySelector("g.key");
  animate(key, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector('circle[cx="10"]');
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
