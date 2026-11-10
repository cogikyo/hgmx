import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelectorAll('path[d^="M7.507"], circle[cx="5.5"]');
  animate(left, { pathLength: [0, 1], opacity: [0, 1], translateX: [-3, 0] }, { duration: τmid });

  const right = el.querySelectorAll('path[d^="M22 10.5"], circle[cx="18.5"]');
  animate(right, { pathLength: [0, 1], opacity: [0, 1], translateX: [3, 0] }, { duration: τmid });

  const body = el.querySelector('path[d^="M17 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const head = el.querySelector('circle[cx="12"]');
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
