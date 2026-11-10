import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const rider = el.querySelector("g.rider");
  animate(rider, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });

  const arrow = el.querySelector("g.arrow");
  animate(arrow, { translateX: [0, 2, -1, 0], translateY: [0, -2, 1, 0] }, { duration: τmid, delay: τmin });
}
