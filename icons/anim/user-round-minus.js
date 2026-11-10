import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M2 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const minus = el.querySelector('path[d="M22 19h-6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
