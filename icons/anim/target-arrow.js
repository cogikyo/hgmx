import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g.arrow");
  animate(arrow, { translateX: [6, -0.5, 0], translateY: [-6, 0.5, 0], opacity: [0, 1, 1] }, { duration: τmid });
  const rings = el.querySelector("g.rings");
  animate(rings, { scale: [1, 0.9, 1] }, { duration: τmin, delay: τmid });
}
