import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const up = el.querySelector("g.up");
  animate(up, { translateX: [0, 2, -0.5, 0], translateY: [0, -2, 0.5, 0] }, { duration: τmid });

  const down = el.querySelector("g.down");
  animate(down, { translateX: [0, 2, -0.5, 0], translateY: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
