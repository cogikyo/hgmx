import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const down = el.querySelector("g.down");
  animate(down, { translateY: [0, 3, -0.5, 0] }, { duration: τmid });
  const up = el.querySelector("g.up");
  animate(up, { translateY: [0, -3, 0.5, 0] }, { duration: τmid, delay: τmin });
}
