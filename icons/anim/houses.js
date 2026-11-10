import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector('path[d^="M15.485"]');
  animate(front, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const back = el.querySelector('path[d^="m12.681"]');
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const door = el.querySelector('path[d^="M7.41"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
