import { animate, stagger } from "../motion.js";

export default function (el) {
  const balls = [[4, 19], [12, 19], [20, 19], [8, 12], [16, 12], [12, 5]].map(([x, y]) => el.querySelector(`circle[cx="${x}"][cy="${y}"]`));
  animate(balls, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
