import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hill = el.querySelector('path[d^="m20 17"]');
  animate(hill, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector("circle");
  animate(sun, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
