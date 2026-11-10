import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const jacket = el.querySelector("g");
  animate(jacket, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });

  const stripes = el.querySelectorAll('path[d="M2 15h4l2-2"], path[d="M22 15h-4l-2-2"], path[d="M6 18h12"]');
  animate(stripes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
