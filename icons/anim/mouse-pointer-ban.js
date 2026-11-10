import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pointer = el.querySelector('path[d^="M2.034"]');
  animate(pointer, { translateX: [0, -2, 1, 0], translateY: [0, -2, 1, 0] }, { duration: τmid });

  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m11.8"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
