import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pin = el.querySelector("g");
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1] }, { duration: τmid });

  const house = el.querySelector('path[d^="M15 22"]');
  animate(house, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmin });

  const door = el.querySelector('path[d="M18 22v-3"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
