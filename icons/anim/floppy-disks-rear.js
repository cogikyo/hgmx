import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g");
  animate(front, { translateX: [-4, 1, 0], translateY: [4, -1, 0] }, { duration: τmax });

  const back = el.querySelector('path[d^="M18 22"]');
  animate(back, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid, delay: τmin });

  const hub = el.querySelector("circle");
  animate(hub, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
