import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M21 12.17"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const play = el.querySelector('path[d^="M15 15"]');
  animate(play, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
